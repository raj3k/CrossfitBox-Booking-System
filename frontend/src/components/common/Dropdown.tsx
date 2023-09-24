import { ReactNode, useEffect, useRef } from "react";
import useToggle from "../../hooks/useToggle";
import Icon from "../Icon";

interface Props {
  trigger?: ReactNode;
  actions?: ReactNode;
  className?: string;
  actionsClassName?: string;
}

const Dropdown: React.FC<Props> = (props: Props) => {
  const { trigger, actions, className, actionsClassName } = props;
  const [dropdownStatus, toggleDropdownStatus] = useToggle(false);
  const dropdownWrapperRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (dropdownStatus) {
      const handleClickOutside = (event: MouseEvent) => {
        if (!dropdownWrapperRef.current?.contains(event.target as Node)) {
          toggleDropdownStatus(false);
        }
      };

      window.addEventListener("click", handleClickOutside, {
        capture: true,
      });
      return () => {
        window.addEventListener("click", handleClickOutside, {
          capture: true,
        });
      };
    }
  }, [dropdownStatus]);

  const handleToggleDropdownStatus = (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
    e.stopPropagation;
    toggleDropdownStatus();
  };

  return (
    <div 
      ref={dropdownWrapperRef} 
      className={`relative flex flex-col justify-start items-start select-none ${className ?? ""}`} 
      onClick={handleToggleDropdownStatus}
    >
      {trigger ? (
        trigger
      ) : (
        <button className="flex flex-row justify-center items-center rounded text-gray-400 cursor-pointer hover:text-gray-500">
          <Icon.MoreVertical className="w-4 h-auto" />
        </button>
      )}
      <div
        className={`w-auto mt-1 absolute top-full right-0 flex flex-col justify-start items-start bg-white z-1 border p-1 rounded-md shadow ${
          actionsClassName ?? ""
        } ${dropdownStatus ? "" : "!hidden"}`}
      >
        {actions}
      </div>
    </div>
  );
};

export default Dropdown;