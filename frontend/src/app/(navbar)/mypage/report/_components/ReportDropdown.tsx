import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { RxHamburgerMenu } from "react-icons/rx";

interface Props {
  items: string[];
}

const ReportDropdown = ({ items }: Props) => {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <button>
          <RxHamburgerMenu />
        </button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-56">
        {items.map((item, index) => {
          return <DropdownMenuItem key={index}>{item}</DropdownMenuItem>;
        })}
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export default ReportDropdown;
