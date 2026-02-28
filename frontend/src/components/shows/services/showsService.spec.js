import showsService from './showsService'
import apiService from "../../../helpers/apiService";
import { vi, describe, it, expect } from 'vitest';

vi.mock('../../../helpers/apiService');

describe('Show Service', () => {

    it('should return all shows', async () => {
        const data = [{
            id: 1,
            name: "n1",
            description: "d1",
            price: 1,
            status: "RUNNING"
        },
            {
                id: 2,
                name: "n2",
                description: "d2",
                price: 4,
                status: "RUNNING"
            }];

        vi.mocked(apiService.get).mockResolvedValue({data: data});
        const shows = await showsService.fetchAll();

        expect(shows).toHaveLength(2);

        expect(shows).toEqual([{
            id: 1,
            name: "n1",
            description: "d1",
            price: 1,
            status: "RUNNING"
        }, {
            id: 2,
            name: "n2",
            description: "d2",
            price: 4,
            status: "RUNNING"
        }]);
    });

    it('should create a show', async () => {
        const payload = {
            name: "Movie Name",
            description: "Movie Description",
            price: 100,
            status: "RUNNING"
        };

        const response = {
            id: 1,
            ...payload
        };

        vi.mocked(apiService.post).mockImplementation((url, data) => {
            if (data === payload) {
                return Promise.resolve({data: response});
            }
            return Promise.reject(new Error('Unexpected payload'));
        });

        const createdShow = await showsService.create(payload);

        expect(createdShow).toEqual({
            id: 1,
            name: "Movie Name",
            description: "Movie Description",
            price: 100,
            status: "RUNNING"
        });
    });
});
