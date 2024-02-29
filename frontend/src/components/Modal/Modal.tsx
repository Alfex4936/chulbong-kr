import CloseIcon from "@mui/icons-material/Close";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import useModalStore from "../../store/useModalStore";
import * as Styled from "./Modal.style";

interface Props {
  children: React.ReactNode;
  setState?: React.Dispatch<React.SetStateAction<boolean>>;
}

const BasicModal = ({ children, setState }: Props) => {
  const modalState = useModalStore();

  return (
    <Styled.ModalWrap>
      <Styled.Modal
        onClick={(e) => {
          e.stopPropagation();
        }}
      >
        <Tooltip title="닫기" arrow disableInteractive>
          <IconButton
            onClick={() => {
              modalState.close();
              if (setState) {
                setState(false);
              }
            }}
            aria-label="delete"
            sx={{
              position: "absolute",
              top: ".4rem",
              right: ".4rem",
            }}
          >
            <CloseIcon />
          </IconButton>
        </Tooltip>
        {children}
      </Styled.Modal>
    </Styled.ModalWrap>
  );
};

export default BasicModal;
