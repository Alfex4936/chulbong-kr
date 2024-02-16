import styled from "@emotion/styled";

export const MapContainer = styled.div`
  width: 100vw;
  height: calc(100vh - 60px);

  margin-top: 60px;
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
