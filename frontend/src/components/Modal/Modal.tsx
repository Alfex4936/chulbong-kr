import CloseIcon from "@mui/icons-material/Close";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import useModalStore from "../../store/useModalStore";
import * as Styled from "./Modal.style";
import { useEffect } from "react";

interface Props {
  exit?: boolean;
  children: React.ReactNode;
  setState?: React.Dispatch<React.SetStateAction<boolean>>;
}

const BasicModal = ({ exit = true, children, setState }: Props) => {
  const modalState = useModalStore();

  useEffect(() => {
    const handleKeyDownClose = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        if (setState) {
          setState(false);
        } else {
          modalState.close();
        }
      }
    };

    window.addEventListener("keydown", handleKeyDownClose);

    return () => {
      window.removeEventListener("keydown", handleKeyDownClose);
    };
  }, []);

  return (
    <Styled.ModalWrap>
      <Styled.Modal
        onClick={(e) => {
          e.stopPropagation();
        }}
      >
        {exit && (
          <Tooltip title="닫기" arrow disableInteractive>
            <IconButton
              onClick={() => {
                if (setState) {
                  setState(false);
                } else {
                  modalState.close();
                }
              }}
              aria-label="delete"
              sx={{
                position: "absolute",
                top: "0",
                right: "0",
              }}
            >
              <CloseIcon />
            </IconButton>
          </Tooltip>
        )}

        {children}
      </Styled.Modal>
    </Styled.ModalWrap>
  );
};

export default BasicModal;
