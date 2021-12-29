export default class UsersService {
  /**
   * @param {AxiosInstance} client
   */
  constructor(client) {
    this.client = client;
    this.apiUrl = "/api/users";
  }

  async get(id = null) {
    const userId = id !== null ? `/${id}` : "";
    const response = await this.client.get(`${this.apiUrl}${userId}`);

    return {
      status: response.status === 200,
      data: id === null ? response.data.collection : response.data,
    };
  }

  async getParticipations(id) {
    const response = await this.client.get(
      `${this.apiUrl}/${id}/participations`
    );

    return {
      status: response.status === 200,
      data: response.data.collection,
    };
  }

  async post(user) {
    const response = await this.client.post(`${this.apiUrl}`, user);

    return {
      status: response.status === 201,
      data: response.data,
    };
  }
}
