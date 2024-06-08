import React, { FunctionComponent, useEffect } from "react";
import { useGenerateCharacter } from "../../hooks/useGenerateCharacter";
import { Link } from "react-router-dom";
import { useCollection } from "../../hooks/useCollection";
import { Authenticated } from "../../components/Authenticated/Authenticated";

export const Dashboard: FunctionComponent = () => {
  const characterCollection = useCollection("characters");
  const questCollection = useCollection("quests");

  return (
    <Authenticated>
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
          <div className="table-container">
            <table className="table is-fullwidth">
              <thead>
                <tr>
                  <th>Title</th>
                </tr>
              </thead>
              <tbody>
                {questCollection.data?.items.map((item) => {
                  return (
                    <tr key={item.id}>
                      <td>
                        <Link to={`/quests/${item.id}`}>{item.title}</Link>
                      </td>
                    </tr>
                  );
                })}
                {!questCollection.data ||
                questCollection.data?.items.length === 0 ? (
                  <tr>
                    <td colSpan={2} className="has-text-centered">
                      No quest yet. <Link to="/quests/new">Create one</Link>
                    </td>
                  </tr>
                ) : null}
              </tbody>
            </table>
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
                {!characterCollection.data ||
                characterCollection.data?.items.length === 0 ? (
                  <tr>
                    <td colSpan={2} className="has-text-centered">
                      No character yet.{" "}
                      <Link to="/characters/new">Create one</Link>
                    </td>
                  </tr>
                ) : null}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </Authenticated>
  );
};
