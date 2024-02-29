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

export const imageWrap = styled.div`
  position: relative;

  width: 90%;

  margin: auto;
  margin-bottom: 2rem;
`;

export const SkeletonImage = styled.div`
  display: inline-block;

  height: 300px;
  width: 90%;

  background: #f6f7f8;
  background-image: linear-gradient(
    90deg,
    #f0f0f0 25%,
    #f7f7f7 50%,
    #f0f0f0 75%
  );

  animation: ${shimmer} 1.2s ease-in-out infinite;

  border-radius: 1rem;
`;

export const SkeletonButtons = styled.div`
  position: absolute;
  bottom: 1rem;
  left: 50%;
  transform: translateX(-50%);

  height: 30px;
  width: 30px;

  background: #f6f7f8;
  background-image: linear-gradient(
    90deg,
    #f0f0f0 25%,
    #f7f7f7 50%,
    #f0f0f0 75%
  );

  animation: ${shimmer} 1.2s ease-in-out infinite;
  border-radius: 50%;
`;

// export const BottomButtons = styled.div``;
