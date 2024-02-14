import styled from "@emotion/styled";

export const HeaderContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;

  position: absolute;
  top: 0;
  left: 0;

  padding: 0 2rem;

  width: 100vw;
  height: 60px;

  background-color: #fff;
  color: #333;

  box-shadow: rgba(50, 50, 93, 0.25) 0px 2px 5px -1px,
    rgba(0, 0, 0, 0.3) 0px 1px 3px -1px;

  z-index: 110;
`;

export const Modal = styled.div`
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);

  padding: 2rem;

  width: 400px;

  background-color: #fff;

  border-radius: 1rem;
`;
