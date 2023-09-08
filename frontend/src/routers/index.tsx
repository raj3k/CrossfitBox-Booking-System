import { createBrowserRouter } from 'react-router-dom';
import App from "../App.tsx";
import SignIn from "../pages/SignIn.tsx";
import SignUp from "../pages/SignUp.tsx";
import Root from "../layouts/Root.tsx";
import Home from "../pages/Home.tsx";

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
    children: [
      {
        path: "auth",
        element: <SignIn />,
      },
      {
        path: "auth/signup",
        element: <SignUp />,
      },
      {
        path: "",
        element: <Root />,
        children: [
          {
            path: "",
            element: <Home />
          },
        ],
      },
    ],
  },
]);

export default router;