class ApiFetcher<T> {
  public async fetchData(
    fetchFunction: (...args: any[]) => Promise<T>,
    ...args: any[]
  ): Promise<{ data?: T; isLoading: boolean; isError: boolean }> {
    let isLoading = true;
    let isError = false;
    let data: T | undefined;

    try {
      data = await fetchFunction(...args);
    } catch (error) {
      console.error("Error fetching data:", error);
      isError = true;
    } finally {
      isLoading = false;
    }

    return { data, isLoading, isError };
  }
}
