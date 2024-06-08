import React, { FunctionComponent, useCallback, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { usePocket } from "../../contexts/pocketbase";
export const NavBar: FunctionComponent = () => {
  const { user, logout } = usePocket();
  const navigate = useNavigate();
  const [isActive, setActive] = useState(false);

  const onNavbarBurgerClick = useCallback(() => {
    setActive((active) => !active);
  }, []);

  const onLogoutClick = useCallback(() => {
    logout();
    navigate("/");
  }, []);

  return (
    <nav className="navbar" role="navigation" aria-label="main navigation">
      <div className="container">
        <div className="navbar-brand">
          <a className="navbar-item" href="/">
            <span className="title">Side</span>
            <span className="title has-text-info">Quest</span>
          </a>

          <a
            role="button"
            className="navbar-burger has-text-black"
            aria-label="menu"
            aria-expanded={isActive ? "true" : "false"}
            onClick={onNavbarBurgerClick}
          >
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
          </a>
        </div>

        <div className={`navbar-menu ${isActive ? "is-active" : ""}`}>
          <div className="navbar-end">
            <div className="navbar-item">
              <div className="buttons">
                {user ? (
                  <button className="button is-light" onClick={onLogoutClick}>
                    <strong>Logout</strong>
                  </button>
                ) : (
                  <Link to="/login" className="button is-light">
                    <strong>Login</strong>
                  </Link>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </nav>
  );
};
