import { cn } from "@/utils/cn";

type Props = {
  status: string;
  className?: string;
};

const StatusBadge = ({ status, className }: Props) => {
  const getStatusText = () => {
    if (status === "PENDING") return "대기중";
    if (status === "APPROVED") return "승인 완료";
    if (status === "DENIED") return "거절";
  };

  const getStatusColor = () => {
    if (status === "PENDING") return "bg-grey-dark-1";
    if (status === "APPROVED") return "bg-green-300";
    if (status === "DENIED") return "bg-red";
  };

  return (
    <div
      className={cn(
        "flex items-center bg-black-light-2 border border-solid border-grey px-2 py-1 rounded-lg text-nowrap",
        className
      )}
    >
      <div className={`w-2 h-2 mr-2 rounded-full ${getStatusColor()}`} />
      <div className="text-xs">{getStatusText()}</div>
    </div>
  );
};

export default StatusBadge;
