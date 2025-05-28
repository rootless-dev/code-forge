import {mountSuspended} from '@nuxt/test-utils/runtime'
import {describe, expect, it} from 'vitest'
import Index from '~/pages/index.vue'

describe('Index', () => {
    it('should should render index screen', async () => {
        const component = await mountSuspended(Index)
        expect(component.html()).contains("#hero")
    });
});