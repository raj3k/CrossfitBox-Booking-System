import { useEffect } from "react";
import {Outlet, useNavigate} from "react-router-dom";
import useUserStore from "../stores/v1/user";

const Root: React.FC = () => {
  const navigate = useNavigate();
  const userStore = useUserStore();
  const currentUser = userStore.getCurrentUser();
  const isInitialized = Boolean(currentUser);

  useEffect(() => {
    if (!currentUser) {
      navigate("/auth", {
        replace: true,
      });
      return;
    }
  }, []);

  return (
    <>
    {isInitialized && (
      <div className="w-full min-h-full bg-white">
        <div className="w-full max-w-6xl mx-auto flex flex-row justify-center items-start sm:px-4">
          <main className="w-full max-w-full flex-grow shrink flex flex-col justify-center items-start">
            <Outlet />
          </main>
        </div>
      </div>
    )}
    </>
  )
}

export default Root;