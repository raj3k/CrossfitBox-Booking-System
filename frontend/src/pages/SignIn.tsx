import {Button, Input} from "@mui/joy";
import { FormEvent, useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { useNavigate } from "react-router-dom";
import { toast } from "react-hot-toast";
import * as api from "../helpers/api";
import useUserStore from "../stores/v1/user";

const SignIn: React.FC = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();
  const userStore = useUserStore();
  const allowConfirm = email.length > 0 && password.length > 0;

  useEffect(() => {
    if (userStore.getCurrentUser()) {
      return navigate("/", {
        replace: true,
      });
    }
  }, []);

  const handleEmailInputChanged = (e: React.ChangeEvent<HTMLInputElement>) => {
    const text = e.target.value as string;
    setEmail(text);
  }

  const handlePasswordInputChanged = (e: React.ChangeEvent<HTMLInputElement>) => {
    const text = e.target.value as string;
    setPassword(text);
  }

  const handleSigninBtnClick = async (e: FormEvent) => {
    e.preventDefault();
    try {
      await api.signIn(email, password);
      
      const user = await userStore.fetchCurrentUser();
      
      if (user) {
        navigate("/", {
          replace: true,
        });
      } else {
        toast.error("Signin failed");
      }
    } catch(error: any) {
      console.error(error);
      toast.error(error.response.data.message);
    }
  }

  return (
    <div className="flex flex-row justify-center items-center w-full h-auto mt-12 sm:mt-24 bg-white">
      <div className="w-80 max-w-full h-full py-4 flex flex-col justify-start items-center">
        <div className="w-full py-4 grow flex flex-col justify-center items-center">
          <form className="w-full mt-4" onSubmit={handleSigninBtnClick}>
            <div className="flex flex-col justify-start items-start w-full gap-4">
              <Input
                className="w-full"
                size="lg"
                type="text"
                placeholder="Email"
                onChange={handleEmailInputChanged}
                required
              />
              <Input
                className="w-full"
                size="lg"
                type="password"
                placeholder="Password"
                onChange={handlePasswordInputChanged}
                required
              />
            </div>
            <div className="flex flex-col justify-center items-center w-full mt-6">
              <Button
                type="submit"
                color="primary"
                disabled={!allowConfirm}
                onClick={handleSigninBtnClick}
              >
                Sign In
              </Button>
            </div>
          </form>
          <p className="w-full mt-4 text-sm">
            <span>{"Don't have an account yet?"}</span>
            <Link to="/auth/signup" className="cursor-pointer ml-2 text-blue-600 hover:underline">
              Sign up
            </Link>
          </p>
        </div>
      </div>
    </div>
  )
}

export default SignIn;