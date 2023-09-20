import { useEffect, useState } from "react";
import {Outlet} from "react-router-dom";
import useUserStore from "./stores/v1/user";

function App() {
  const userStore = useUserStore();
  const [loading, setLoading] = useState(true);
  useEffect(() => {
    const initialState = async () => {
      try {
        await userStore.fetchCurrentUser();
      } catch (error) {
        // do nothing
      }
      setLoading(false);
    };

    initialState();
  }, []);

  return !loading ? (
    <>
      <Outlet />
    </>
  ) : (
    <></>
  )
}

export default App
