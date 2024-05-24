import { cn } from "@/utils/cn";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";

interface Props {
  icon: string | React.ReactElement;
  tooltipText: string;
  onClick: VoidFunction;
  className: string;
}

const MapButtons = ({ icon, tooltipText, onClick, className }: Props) => {
  return (
    <TooltipProvider>
      <Tooltip delayDuration={10}>
        <TooltipTrigger
          className={cn(
            "absolute bg-black-light-2 p-[6px] text-xl rounded-sm z-50 hover:bg-grey-dark-1",
            className
          )}
          onClick={onClick}
        >
          {icon}
        </TooltipTrigger>
        <TooltipContent className="bg-black-light-2 text-grey">
          <p>{tooltipText}</p>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  );
};

export default MapButtons;
