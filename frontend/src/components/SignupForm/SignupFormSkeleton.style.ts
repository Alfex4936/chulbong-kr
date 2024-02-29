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

export const TitleSkeleton = styled.div`
  display: inline-block;

  height: 40px;
  width: 100px;

  background: #f6f7f8;
  background-image: linear-gradient(
    90deg,
    #f0f0f0 25%,
    #f7f7f7 50%,
    #f0f0f0 75%
  );

  animation: ${shimmer} 1.2s ease-in-out infinite;

  border-radius: 0.5rem;

  margin: 1rem auto;
`;

export const InputSkeleton = styled.div`
  display: inline-block;

  height: 30px;
  width: 100%;

  background: #f6f7f8;
  background-image: linear-gradient(
    90deg,
    #f0f0f0 25%,
    #f7f7f7 50%,
    #f0f0f0 75%
  );

  animation: ${shimmer} 1.2s ease-in-out infinite;

  margin-bottom: 0.5rem;
`;

export const ButtonSkeleton = styled.div`
  display: inline-block;

  height: 35px;
  width: 95%;

  background: #f6f7f8;
  background-image: linear-gradient(
    90deg,
    #f0f0f0 25%,
    #f7f7f7 50%,
    #f0f0f0 75%
  );

  animation: ${shimmer} 1.2s ease-in-out infinite;

  border-radius: 0.5rem;

  margin: 1rem 0;
`;
