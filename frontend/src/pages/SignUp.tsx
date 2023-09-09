import {Button, Input} from "@mui/joy";
import {FormEvent, useState} from "react";
import * as api from "../helpers/api";

const SignUp: React.FC = () => {
  const [email, setEmail] = useState("");
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [password, setPassword] = useState("");

  const handleEmailInputChanged = (e: React.ChangeEvent<HTMLInputElement>) => {
    const text = e.target.value as string;
    setEmail(text);
  };

  const handleFirstNameInputChanged = (e: React.ChangeEvent<HTMLInputElement>) => {
    const text = e.target.value as string;
    setFirstName(text);
  };

  const handleLastNameInputChanged = (e: React.ChangeEvent<HTMLInputElement>) => {
    const text = e.target.value as string;
    setLastName(text);
  };

  const handlePasswordInputChanged = (e: React.ChangeEvent<HTMLInputElement>) => {
    const text = e.target.value as string;
    setPassword(text);
  };

  const handleSignUpBtnClick = async (e: FormEvent)=> {
    e.preventDefault();
    try {
      await api.signup(email, firstName, lastName, password);
    } catch (error: any) {
      console.error(error);
    }
  }

  return (
    <div className="flex flex-row justify-center items-center w-full h-auto mt-12 sm:mt-24 bg-white">
      <div className="w-80 max-w-full h-full py-4 flex flex-col justify-start items-center">
        <div className="w-full py-4 grow flex flex-col justify-center items-center">
          <p className="w-full text-2xl mt-6">Create your account</p>
          <form className="w-full mt-4" onSubmit={handleSignUpBtnClick}>
            <div className="flex flex-col justify-start items-start w-full gap-4 py-4">
              <Input
                className="w-full"
                size="lg"
                type="email"
                placeholder="Email"
                value={email}
                onChange={handleEmailInputChanged}
                required
              />
              <Input
                className="w-full"
                size="lg"
                type="text"
                placeholder="First Name"
                value={firstName}
                onChange={handleFirstNameInputChanged}
                required
              />
              <Input
                className="w-full"
                size="lg"
                type="text"
                placeholder="Last Name"
                value={lastName}
                onChange={handleLastNameInputChanged}
                required
              />
              <Input
                className="w-full"
                size="lg"
                type="password"
                placeholder="Password"
                value={password}
                onChange={handlePasswordInputChanged}
                required
              />
            </div>
            <div className="flex flex-col justify-center items-center w-full mt-6">
              <Button
                className="w-full"
                type="submit"
                color="primary"
              >
                Sign Up
              </Button>
            </div>
          </form>
        </div>
      </div>
    </div>
  )
}

export default SignUp;