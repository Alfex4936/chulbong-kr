import styled from "@emotion/styled";

export const BlackContainer = styled.div`
  position: absolute;
  top: 0;
  left: 0;

  width: 100%;
  height: 100%;

  z-index: 100;

  background-color: ${({ bg }: { bg: "transparent" | "black" }) =>
    bg === "transparent" ? "transparent" : "rgba(0, 0, 0, 0.6)"};
`;

export const ChildContainer = styled.div`
  position: absolute;
  top: 50%;
  left: 50%;

  transform: translate(-50%, -50%);
`;
