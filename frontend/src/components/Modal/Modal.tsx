import CloseIcon from "@mui/icons-material/Close";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import useModalStore from "../../store/useModalStore";
import * as Styled from "./Modal.style";

interface Props {
  exit?: boolean;
  children: React.ReactNode;
  setState?: React.Dispatch<React.SetStateAction<boolean>>;
}

const BasicModal = ({ exit = true, children, setState }: Props) => {
  const modalState = useModalStore();

  const navigate = useNavigate();
  const query = new URLSearchParams(location.search);
  const sharedMarker = query.get("d");
  const sharedMarkerLat = query.get("la");
  const sharedMarkerLng = query.get("lo");

  const modalRef = useRef(null);

  useEffect(() => {
    const handleClose = () => {
      if (setState) {
        setState(false);
      } else {
        modalState.close();
      }
    };

    const handleKeyDownClose = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        if (sharedMarker && sharedMarkerLat && sharedMarkerLng) {
          navigate("/");
        }
        handleClose();
      }
    };

    window.addEventListener("keydown", handleKeyDownClose);

    return () => {
      window.removeEventListener("keydown", handleKeyDownClose);
    };
  }, []);

  return (
    <Styled.ModalWrap
      onClick={(e) => {
        e.stopPropagation();
        if (sharedMarker && sharedMarkerLat && sharedMarkerLng) {
          navigate("/");
        }

        if (setState) {
          setState(false);
        } else {
          modalState.close();
        }
      }}
    >
      <Styled.Modal
        ref={modalRef}
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
