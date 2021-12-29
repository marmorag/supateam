export default class AuthenticationService {
  /**
   * @param {AxiosInstance} client
   */
  constructor(client) {
    this.client = client;
  }

  async authenticate({ identity }) {
    const response = await this.client.post("/api/auth/login", {
      Identity: identity,
    });

    return {
      status: response.status === 200,
      data: response.data.token ?? null,
    };
  }
}
