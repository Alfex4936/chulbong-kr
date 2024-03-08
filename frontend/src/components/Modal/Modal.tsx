import CloseIcon from "@mui/icons-material/Close";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import useModalStore from "../../store/useModalStore";
import * as Styled from "./Modal.style";

interface Props {
  exit?: boolean;
  children: React.ReactNode;
  setState?: React.Dispatch<React.SetStateAction<boolean>>;
}

const BasicModal = ({ exit = true, children, setState }: Props) => {
  const modalState = useModalStore();

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
                modalState.close();
                if (setState) {
                  setState(false);
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
