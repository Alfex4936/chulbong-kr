import styled from "@emotion/styled";
import { keyframes } from "@emotion/react";

const shimmer = keyframes`
  0% {
    background-position: -468px 0;
  }
  100% {
    background-position: 468px 0;
  }
`;

export const ListSkeleton = styled.div`
  display: flex;
  align-items: center;

  margin: 0 auto 1rem auto;

  padding: 1rem;

  height: 57px;
  width: 100%;

  background: #f6f7f8;
  background-image: linear-gradient(
    90deg,
    #f0f0f0 25%,
    #f7f7f7 50%,
    #f0f0f0 75%
  );

  animation: ${shimmer} 1.2s ease-in-out infinite;
  border-radius: 0.4rem;
`;
