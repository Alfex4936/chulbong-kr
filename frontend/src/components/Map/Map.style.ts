import styled from "@emotion/styled";

export const MapContainer = styled.div`
  width: 100vw;
  height: 100vh;
`;

export const AlertContainer = styled.div`
  border: 1px solid red;

  position: absolute;
  left: 0;
  bottom: 1rem;

  width: 50%;
  height: 80px;

  background-color: #333;

  z-index: 1;
`;

export const ExitButton = styled.div`
  position: absolute;
  top: 0;
  right: 0;

  border-radius: 50%;

  width: 25px;
  height: 25px;

  &:hover {
    color: rgb(230, 103, 103);
  }
`;

export const DeleteUserButtonsWrap = styled.div`
  display: flex;

  & button {
    margin: 1rem 1rem 0 1rem;
  }
`;

export const ChangePasswordButtonsWrap = styled.div`
  display: flex;

  & button {
    margin: 1rem 0.5rem 0 0.5rem;
  }
`;

export const ErrorBox = styled.div`
  text-align: left;

  font-size: 0.7rem;

  color: red;
`;
