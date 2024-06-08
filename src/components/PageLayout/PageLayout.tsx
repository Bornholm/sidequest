import React, { FunctionComponent, Fragment } from "react";
import { Outlet } from "react-router";
import { NavBar } from "../NavBar/NavBar";

export const PageLayout: FunctionComponent = () => {
  return (
    <Fragment>
      <NavBar />
      <section className="section">
        <div className="container">
          <Outlet />
        </div>
      </section>
    </Fragment>
  );
};
