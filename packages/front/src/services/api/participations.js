export default class ParticipationsService {
  /**
   * @param {AxiosInstance} client
   */
  constructor(client) {
    this.client = client;
    this.apiUrl = "/api/participations";
  }

  async get(id = null) {
    const participationId = id !== null ? `/${id}` : "";
    const response = await this.client.get(`${this.apiUrl}${participationId}`);

    return {
      status: response.status === 200,
      data: id === null ? response.data.collection : response.data,
    };
  }

  async put(participation) {
    console.log(this.client);
    const response = await this.client.put(
      `${this.apiUrl}/${participation.id}`,
      participation
    );

    return {
      status: response.status === 200,
      data: response.data,
    };
  }

  async post(participation) {
    const response = await this.client.post(`${this.apiUrl}`, participation);

    return {
      status: response.status === 201,
      data: response.data,
    };
  }
}
