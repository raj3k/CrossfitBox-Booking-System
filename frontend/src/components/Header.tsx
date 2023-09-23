import { Link } from "react-router-dom";

const Header: React.FC = () => {
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
          </div>
        </div>
      </div>
			</>
    );
}

export default Header;