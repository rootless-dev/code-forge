import {describe, expect, it} from "vitest";
import {renderSuspended} from "@nuxt/test-utils/runtime";
import IndexHeaderBar from "~/components/ui/headers/IndexHeaderBar/index.vue";
import {screen} from "@testing-library/dom";

describe('IndexHeaderBar', () => {
    const mockSections = [
        {label: 'Home', to: '/'},
        {label: 'About', to: '/about'}
    ]

    it('should render components', async () => {
        await renderSuspended(IndexHeaderBar, {
            props: {
                sections: mockSections
            }
        })

        const logo = screen.getByRole('img')
        expect(logo).toBeTruthy()

        expect(screen.getByText('Sign In')).toBeTruthy()
        expect(screen.getByText('Sign Up')).toBeTruthy()

        expect(screen.getByText('Home')).toBeTruthy()
        expect(screen.getByText('About')).toBeTruthy()
    });
})