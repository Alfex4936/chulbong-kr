import styled from "@emotion/styled";

export const Container = styled.div`
  z-index: 900;
`;

export const RoadViewContainer = styled.div`
  position: fixed;
  top: 0;
  left: 0;

  width: 100%;
  height: 100%;
`;

export const Exit = styled.div`
  position: absolute;
  top: 5px;
  left: 3px;

  background-color: rgba(0, 0, 0, 0.5);

  border-radius: 0.5rem;

  font-size: 1rem;
  z-index: 3;
`;
