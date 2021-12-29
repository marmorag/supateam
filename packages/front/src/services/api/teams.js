export default class TeamsService {
  /**
   * @param {AxiosInstance} client
   */
  constructor(client) {
    this.client = client;
    this.apiUrl = "/api/teams";
  }

  async get(id = null) {
    const eventId = id !== null ? `/${id}` : "";
    const response = await this.client.get(`${this.apiUrl}${eventId}`);

    return {
      status: response.status === 200,
      data: id === null ? response.data.collection : response.data,
    };
  }

  async post(event) {
    const response = await this.client.post(`${this.apiUrl}`, event);

    return {
      status: response.status === 201,
      data: response.data,
    };
  }
}
