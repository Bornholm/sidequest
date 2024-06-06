import React, { FunctionComponent, useEffect, useCallback } from "react";
import { usePocket } from "../../contexts/pocketbase";
import { useNavigate } from "react-router";

export const Home: FunctionComponent = () => {
  const { user } = usePocket();
  const navigate = useNavigate();

  useEffect(() => {
    if (!user) return;
    navigate("/dashboard");
  }, [user]);

  return (
    <section className="hero is-medium">
      <div className="hero-body">
        <p className="title">
          Become the MJ-hero of your tabletop RPG friends !
        </p>
        <p className="subtitle">
          Easily use AI to generate quest plots, characters and much, much
          more...
        </p>
      </div>
    </section>
  );
};
