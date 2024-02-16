import * as Styled from "./Modal.style";
import useModalStore from "../../store/useModalStore";

interface Props {
  children: React.ReactNode;
  setState?: React.Dispatch<React.SetStateAction<boolean>>;
}

const BasicModal = ({ children, setState }: Props) => {
  const modalState = useModalStore();

  return (
    <Styled.ModalWrap
      onClick={() => {
        modalState.close();
        if (setState) {
          setState(false);
        }
      }}
    >
      <Styled.Modal
        onClick={(e) => {
          e.stopPropagation();
        }}
      >
        {children}
      </Styled.Modal>
    </Styled.ModalWrap>
  );
};

export default BasicModal;
