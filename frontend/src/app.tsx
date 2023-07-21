import "./App.css";
import LaunchScreen from "./screens/LaunchScreen";
import LoginScreen from "./screens/LoginScreen";
import HomeScreen from "./screens/HomeScreen";
import { theme } from "./constant/config";
import {
  RouterProvider,
  createBrowserRouter,
  createHashRouter,
} from "react-router-dom";
import ScanScreen from "./screens/ScanScreen";
import { useEffect, useState } from "react";
import { CheckLoginStatus } from "../wailsjs/go/service/Auth";
import { EventsOn } from "../wailsjs/runtime/runtime";

function App() {
  const [authenticated, setAuthenticated] = useState(false);
  const [ready, setReady] = useState(false);

  useEffect(() => {
    CheckLoginStatus().then((state) => setAuthenticated(state.Authenticated));

    EventsOn("auth.state", (state) => {
      setAuthenticated(state.Authenticated);
    });

    setTimeout(() => setReady(true), 2000);
  }, []);

  return (
    <div
      id="app"
      className="App"
      style={{
        backgroundColor: theme.palette.secondary.main,
      }}
    >
      <RouterProvider
        router={createHashRouter([
          {
            path: "/",
            element: ready ? (
              authenticated ? (
                <HomeScreen />
              ) : (
                <LoginScreen />
              )
            ) : (
              <LaunchScreen />
            ),
          },
          {
            path: "/scan",
            element: <ScanScreen />,
          },
        ])}
      />
    </div>
  );
}

export default App;
