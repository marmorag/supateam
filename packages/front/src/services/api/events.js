export const KIND_GRAND_PRIX = "Grand Prix";
export const KIND_EQUIPE = "Equipe";
export const KIND_ENTRAINEMENT = "Entrainement";

export default class EventsService {
  /**
   * @param {AxiosInstance} client
   */
  constructor(client) {
    this.client = client;
    this.apiUrl = "/api/events";
  }

  async get(id = null) {
    const eventId = id !== null ? `/${id}` : "";
    const response = await this.client.get(`${this.apiUrl}${eventId}`);

    return {
      status: response.status === 200,
      data: id === null ? response.data.collection : response.data,
    };
  }

  async getParticipation(id) {
    const response = await this.client.get(
      `${this.apiUrl}/${id}/participations`,
      { params: { format: "long" } }
    );

    return {
      status: response.status === 200,
      data: response.data.collection,
    };
  }

  async post(event) {
    const response = await this.client.post(`${this.apiUrl}`, event);

    return {
      status: response.status === 201,
      data: response.data,
    };
  }

  async put(event) {
    const response = await this.client.put(`${this.apiUrl}/${event.Id}`, event);

    return {
      status: response.status === 200,
      data: response.data,
    };
  }
}
