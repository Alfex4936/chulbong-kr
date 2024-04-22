import Link from "next/link";
import { ArrowLeftIcon } from "../icons/ArrowIcons";

type Props = {
  url: string;
  text?: string;
};

const PrevHeader = ({ url, text }: Props) => {
  return (
    <div className="flex items-center h-14">
      <Link
        href={url}
        className="flex justify-center items-center w-10 h-10 mr-2"
      >
        <ArrowLeftIcon />
      </Link>
      {text && <span>{text}</span>}
    </div>
  );
};

export default PrevHeader;
