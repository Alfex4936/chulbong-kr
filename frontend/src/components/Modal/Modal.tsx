import * as Styled from "./Modal.style";
import useModalStore from "../../store/useModalStore";

interface Props {
  children: React.ReactNode;
}

const BasicModal = ({ children }: Props) => {
  const modalState = useModalStore();

  return (
    <Styled.ModalWrap
      onClick={() => {
        modalState.closeLogin();
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
