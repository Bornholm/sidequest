import React, {
  FunctionComponent,
  Fragment,
  PropsWithChildren,
  useEffect,
} from "react";
import { Outlet, useNavigate } from "react-router";
import { NavBar } from "../NavBar/NavBar";
import { usePocket } from "../../contexts/pocketbase";

export interface AuthenticatedProps extends PropsWithChildren {}

export const Authenticated: FunctionComponent<AuthenticatedProps> = ({
  children,
}) => {
  const { user } = usePocket();
  const navigate = useNavigate();

  useEffect(() => {
    if (user) return;
    navigate("/login");
  }, [user]);

  if (!user) return null;

  return <Fragment>{children}</Fragment>;
};
