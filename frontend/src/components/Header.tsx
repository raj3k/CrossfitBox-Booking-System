import { Link } from "react-router-dom";
import Dropdown from "./common/Dropdown";
import { Avatar } from "@mui/joy";
import Icon from "./Icon";
import useUserStore from "../stores/v1/user";
import * as api from "../helpers/api";

const Header: React.FC = () => {
  const currentUser = useUserStore().getCurrentUser();
  
  const handleSignOutButtonClick = async () => {
    await api.signOut();
    window.location.href = "/auth";
  };

    return (
			<>
				<div className="w-full bg-gray-50 border-b border-b-gray-200">
        <div className="w-full max-w-6xl mx-auto px-3 md:px-12 py-5 flex flex-row justify-between items-center">
          <div className="flex flex-row justify-start items-center shrink mr-2">
            <Link to="/" className="text-lg cursor-pointer flex flex-row justify-start items-center">
              CrossBoxFit
            </Link>
          </div>
          <div className="relative flex-shrink-0">
            <Dropdown
              trigger={
                <button className="flex flex-row justify-end items-center cursor-pointer">
                  <Avatar size="sm" variant="plain" />
                  <span>{currentUser.first_name}</span>
                  <Icon.ChevronDown className="ml-2 w-5 h-auto text-gray-600" />
                </button>
              }
              actionsClassName="!w-32"
              actions={
                <>
                  <Link
                    to="/setting"
                    className="w-full px-2 flex flex-row justify-start items-center text-left leading-8 cursor-pointer rounded hover:bg-gray-100 disabled:cursor-not-allowed disabled:bg-gray-100 disabled:opacity-60"
                  >
                    <Icon.Settings className="w-4 h-auto mr-2" /> Setting
                  </Link>
                  <button
                    className="w-full px-2 flex flex-row justify-start items-center text-left leading-8 cursor-pointer rounded hover:bg-gray-100 disabled:cursor-not-allowed disabled:bg-gray-100 disabled:opacity-60"
                    onClick={() => handleSignOutButtonClick()}
                  >
                    <Icon.LogOut className="w-4 h-auto mr-2" /> Sign out
                  </button>
                </>
              }
            ></Dropdown>
          </div>
        </div>
      </div>
			</>
    );
}

export default Header;