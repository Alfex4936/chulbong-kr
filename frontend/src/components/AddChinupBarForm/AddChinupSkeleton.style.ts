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

export const UploadSkeleton = styled.div`
  display: inline-block;

  height: 150px;
  width: 150px;

  background: #f6f7f8;
  background-image: linear-gradient(
    90deg,
    #f0f0f0 25%,
    #f7f7f7 50%,
    #f0f0f0 75%
  );

  animation: ${shimmer} 1.2s ease-in-out infinite;

  border-radius: 0.5rem;

  margin: auto;
  margin-bottom: 2rem;
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
`;

export const ButtonSkeleton = styled.div`
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

  border-radius: 0.5rem;

  margin-top: 1rem;
`;
