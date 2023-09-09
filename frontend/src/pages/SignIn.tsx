import {Button, Input} from "@mui/joy";

const SignIn: React.FC = () => {
  return (
    <div className="flex flex-row justify-center items-center w-full h-auto mt-12 sm:mt-24 bg-white">
      <div className="w-80 max-w-full h-full py-4 flex flex-col justify-start items-center">
        <div className="w-full py-4 grow flex flex-col justify-center items-center">
          <form className="w-full mt-4">
            <div className="flex flex-col justify-start items-start w-full gap-4">
              <Input
                className="w-full"
                size="lg"
                type="text"
                placeholder="Email"
                required
              />
              <Input
                className="w-full"
                size="lg"
                type="password"
                placeholder="Password"
                required
              />
            </div>
            <div className="flex flex-col justify-center items-center w-full mt-6">
              <Button>
                Sign In
              </Button>
            </div>
          </form>
        </div>
      </div>
    </div>
  )
}

export default SignIn;