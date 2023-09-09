import {Outlet} from "react-router-dom";

const Root: React.FC = () => {
  return (
    <div className="w-full min-h-full bg-white">
      <div className="w-full max-w-6xl mx-auto flex flex-row justify-center items-start sm:px-4">
        <main className="w-full max-w-full flex-grow shrink flex flex-col justify-center items-start">
          <Outlet />
        </main>
      </div>
    </div>
  )
}

export default Root;