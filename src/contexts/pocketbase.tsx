import React, {
  createContext,
  useContext,
  useCallback,
  useState,
  useEffect,
  useMemo,
  FunctionComponent,
  PropsWithChildren,
} from "react";
import PocketBase, {
  AuthModel,
  RecordAuthResponse,
  RecordModel,
} from "pocketbase";
import { useInterval } from "usehooks-ts";
import { jwtDecode } from "jwt-decode";
import ms from "ms";

const fiveMinutesInMs = ms("5 minutes");
const twoMinutesInMs = ms("2 minutes");

export interface Context {
  pb?: PocketBase;
  user?: AuthModel;
  token?: string;
  register: (email: string, password: string) => Promise<RecordModel>;
  login: (provider: string) => Promise<RecordAuthResponse<RecordModel>>;
  logout: () => void;
}

const PocketContext = createContext<Context>({
  register: () => Promise.reject("not yet initialized"),
  login: () => Promise.reject("not yet initialized"),
  logout: () => {},
});

export interface PocketProviderProps extends PropsWithChildren {
  baseUrl: string;
}

export const PocketProvider: FunctionComponent<PocketProviderProps> = ({
  children,
  baseUrl,
}) => {
  const pb: PocketBase = useMemo(() => new PocketBase(baseUrl), [baseUrl]);

  const [token, setToken] = useState(pb.authStore.token);
  const [user, setUser] = useState(pb.authStore.model);

  useEffect(() => {
    return pb.authStore.onChange((token, model) => {
      setToken(token);
      setUser(model);
    });
  }, []);

  const register = useCallback(async (email, password) => {
    return await pb
      .collection("users")
      .create({ email, password, passwordConfirm: password });
  }, []);

  const login = useCallback(async (provider: string) => {
    return await pb.collection("users").authWithOAuth2({ provider });
  }, []);

  const logout = useCallback(() => {
    pb.authStore.clear();
  }, []);

  const refreshSession = useCallback(async () => {
    if (!pb.authStore.isValid) return;
    const decoded = jwtDecode(token);
    const tokenExpiration = decoded.exp;
    if (tokenExpiration !== undefined) {
      const expirationWithBuffer = (decoded.exp + fiveMinutesInMs) / 1000;
      if (tokenExpiration < expirationWithBuffer) {
        await pb.collection("users").authRefresh();
      }
    }
  }, [token]);

  useInterval(refreshSession, token ? twoMinutesInMs : null);

  return (
    <PocketContext.Provider
      value={{ register, login, logout, user, token, pb }}
    >
      {children}
    </PocketContext.Provider>
  );
};

export const usePocket = () => useContext(PocketContext);
