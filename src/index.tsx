import React, { FunctionComponent } from "react";
import ReactDOM from "react-dom";
import { PocketProvider } from "./contexts/pocketbase";
import { Home } from "./pages/Home/Home";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import { PageLayout } from "./components/PageLayout/PageLayout";
import { Login } from "./pages/Login/Login";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Dashboard } from "./pages/Dashboard/Dashboard";
import { EditCharacter } from "./pages/Character/EditCharacter";
import { EditQuest } from "./pages/Quest/EditQuest";

const queryClient = new QueryClient();

const App: FunctionComponent = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <PocketProvider baseUrl="/">
        <BrowserRouter>
          <Routes>
            <Route element={<PageLayout />}>
              <Route path="/" index element={<Home />} />
              <Route path="/login" element={<Login />} />
              <Route path="/dashboard" element={<Dashboard />} />
              <Route
                path="/characters/:characterId"
                element={<EditCharacter />}
              />
              <Route path="/quests/:questId" element={<EditQuest />} />
            </Route>
          </Routes>
        </BrowserRouter>
      </PocketProvider>
    </QueryClientProvider>
  );
};

ReactDOM.render(<App />, document.getElementById("root"));
