import {describe, expect, it} from 'vitest'
import {renderSuspended} from "@nuxt/test-utils/runtime";
import Index from "~/components/ui/footers/IndexFooter/index.vue";

describe('index footer', () => {
    it('should render correct text footer', async () => {
        const component = await renderSuspended(Index)
        const current_year = new Date().getFullYear()
        expect(component.html()).contains('All rights reserved. Alves Solutions 2022-')
        expect(component.html()).contains(current_year.toString())
    });
});