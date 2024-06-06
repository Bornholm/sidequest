import React, { FunctionComponent } from "react";
import { usePocket } from "../../contexts/pocketbase";

export const Dashboard: FunctionComponent = () => {
  const { user } = usePocket();
  return (
    <section className="hero is-medium">
      <h1 className="title">
        Welcome <strong>{user?.username ?? "unknown"}</strong> !
      </h1>
    </section>
  );
};
