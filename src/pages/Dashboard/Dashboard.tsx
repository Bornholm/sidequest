import React, { FunctionComponent, useEffect } from "react";
import { useGenerateCharacter } from "../../hooks/useGenerateCharacter";
import { Link } from "react-router-dom";
import { useCollection } from "../../hooks/useCollection";

export const Dashboard: FunctionComponent = () => {
  const characterCollection = useCollection("characters");

  return (
    <div className="columns">
      <div className="column">
        <div className="level is-mobile">
          <div className="level-left">
            <h2 className="title level-item">My quests</h2>
          </div>
          <div className="level-right">
            <Link to="/quests/new" className="button is-primary">
              <strong>+</strong>
            </Link>
          </div>
        </div>
      </div>
      <div className="column">
        <div className="level is-mobile">
          <div className="level-left">
            <h2 className="title">My characters</h2>
          </div>
          <div className="level-right">
            <Link to="/characters/new" className="button is-primary">
              <strong>+</strong>
            </Link>
          </div>
        </div>
        <div className="table-container">
          <table className="table is-fullwidth">
            <thead>
              <tr>
                <th>Name</th>
                <th>Race</th>
              </tr>
            </thead>
            <tbody>
              {characterCollection.data?.items.map((item) => {
                return (
                  <tr key={item.id}>
                    <td>
                      <Link to={`/characters/${item.id}`}>{item.name}</Link>
                    </td>
                    <td>{item.race}</td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};
