import styled from "@emotion/styled";
import { keyframes } from "@emotion/react";

const fadeIn = keyframes`
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
`;

const slideIn = keyframes`
  from {
    transform: translate(-50%, -60%);
  }
  to {
    transform: translate(-50%, -50%);
  }
`;

export const ModalWrap = styled.div`
  position: absolute;
  top: 0;
  left: 0;

  width: 100vw;
  height: 100vh;

  background-color: rgba(0, 0, 0, 0.6);

  z-index: 1000;

  animation: ${fadeIn} 0.3s ease-out;
`;

export const Modal = styled.div`
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);

  padding: 2rem;

  max-width: 540px;
  min-width: 300px;
  width: 90%;

  max-height: 600px;

  background-color: #fff;

  border-radius: 1rem;

  animation: ${fadeIn} 0.3s ease-out, ${slideIn} 0.3s ease-out;

  box-shadow: rgba(0, 0, 0, 0.1) 0px 4px 12px;
`;
