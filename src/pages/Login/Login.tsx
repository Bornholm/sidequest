import React, {
  FunctionComponent,
  useEffect,
  useCallback,
  Fragment,
  MouseEvent,
  useState,
} from "react";
import { usePocket } from "../../contexts/pocketbase";
import { useAuthMethods } from "../../hooks/useAuthMethods";
import { useNavigate } from "react-router";

export const Login: FunctionComponent = () => {
  const { login, user } = usePocket();
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();
  const authMethodsQuery = useAuthMethods();

  const onProviderClick = useCallback((evt: MouseEvent<HTMLButtonElement>) => {
    const provider = evt.currentTarget.dataset.provider;
    if (!provider) return;
    setError(null);
    login(provider).catch((err) => setError(err.message));
  }, []);

  useEffect(() => {
    if (!user) return;
    navigate("/dashboard");
  }, [user]);

  return (
    <Fragment>
      {error ? (
        <div className="message is-danger">
          <div className="message-body">{error}</div>
        </div>
      ) : null}
      <h1 className="title">Select your identity provider</h1>
      <div className="grid">
        {authMethodsQuery.data?.authProviders.map((provider) => {
          return (
            <div key={provider.name} className="cell">
              <button
                onClick={onProviderClick}
                data-provider={provider.name}
                className="button is-large is-fullwidth"
              >
                <strong>{provider.displayName}</strong>
              </button>
            </div>
          );
        })}
      </div>
    </Fragment>
  );
};
