import axios from "axios";
import apiService from "./apiService";
import { vi, describe, it, expect, beforeEach } from 'vitest';

vi.mock("axios");

describe('Api service', () => {
    let internalServerErrorResponse;
    let otherErrorResponse;

    beforeEach(() => {
        window.location.assign = vi.fn();

        internalServerErrorResponse = {
            response: {
                status: 500
            }
        };

        otherErrorResponse = {
            response: {
                status: 400
            }
        };
    });

    it('Should handle internal server error for get call', async () => {
        vi.mocked(axios.get).mockRejectedValue(internalServerErrorResponse);

        await apiService.get(expect.any(String));
        expect(window.location.assign).toBeCalledTimes(1);
        expect(window.location.assign).toBeCalledWith("/error");
    });

    it('Should handle internal server error for post call', async () => {
        vi.mocked(axios.post).mockRejectedValue(internalServerErrorResponse);

        await apiService.post(expect.any(String), expect.any(Object));
        expect(window.location.assign).toBeCalledTimes(1);
        expect(window.location.assign).toBeCalledWith("/error");
    });

    it('Should rethrow error from get call if not internal server error', async () => {
        vi.mocked(axios.get).mockRejectedValue(otherErrorResponse);

        try {
            await apiService.get(expect.any(String));
            fail("Error not rethrown");
        } catch (error) {
            expect(error).toEqual(otherErrorResponse);
        }
    });

    it('Should rethrow error from get call if not internal server error', async () => {
        vi.mocked(axios.get).mockRejectedValue(otherErrorResponse);

        try {
            await apiService.get(expect.any(String), expect.any(Object));
            fail("Error not rethrown");
        } catch (error) {
            expect(error).toEqual(otherErrorResponse);
        }
    });

    it('Should return promise without handling error for post call', async () => {
        vi.mocked(axios.post).mockRejectedValue(internalServerErrorResponse);

        try {
            await apiService.postWithoutErrorHandling(expect.any(String), expect.any(Object));
        } catch (error) {
            expect(window.location.assign).not.toHaveBeenCalled();
        }
    });
});
