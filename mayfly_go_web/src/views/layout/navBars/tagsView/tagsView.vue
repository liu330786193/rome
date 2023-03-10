<template>
    <div class="layout-navbars-tagsview" :class="{ 'layout-navbars-tagsview-shadow': getThemeConfig.layout === 'classic' }">
        <el-scrollbar ref="scrollbarRef" @wheel.prevent="onHandleScroll">
            <ul class="layout-navbars-tagsview-ul" :class="setTagsStyle" ref="tagsUlRef">
                <li
                    v-for="(v, k) in tagsViewList"
                    :key="k"
                    class="layout-navbars-tagsview-ul-li"
                    :data-name="v.name"
                    :class="{ 'is-active': isActive(v) }"
                    @contextmenu.prevent="onContextmenu(v, $event)"
                    @click="onTagsClick(v, k)"
                    :ref="
                        (el) => {
                            if (el) tagsRefs[k] = el;
                        }
                    "
                >
                    <i class="iconfont icon-webicon318 layout-navbars-tagsview-ul-li-iconfont font14" v-if="isActive(v)"></i>
                    <SvgIcon
                        :name="v.meta.icon"
                        class="layout-navbars-tagsview-ul-li-iconfont"
                        v-if="!isActive(v) && getThemeConfig.isTagsviewIcon"
                    />
                    <span>{{ v.meta.title }}</span>
                    <template v-if="isActive(v)">
                        <SvgIcon
                            name="RefreshRight"
                            class="ml5 layout-navbars-tagsview-ul-li-refresh"
                            @click.stop="refreshCurrentTagsView($route.fullPath)"
                        />
                        <SvgIcon
                            name="Close"
                            class="layout-navbars-tagsview-ul-li-icon layout-icon-active"
                            v-if="!v.meta.isAffix"
                            @click.stop="closeCurrentTagsView(getThemeConfig.isShareTagsView ? v.path : v.path)"
                        />
                    </template>
                    <SvgIcon
                        name="Close"
                        class="layout-navbars-tagsview-ul-li-icon layout-icon-three"
                        v-if="!v.meta.isAffix"
                        @click.stop="closeCurrentTagsView(getThemeConfig.isShareTagsView ? v.path : v.path)"
                    />
                </li>
            </ul>
        </el-scrollbar>
        <Contextmenu :dropdown="dropdown" ref="contextmenuRef" @currentContextmenuClick="onCurrentContextmenuClick" />
    </div>
</template>

<script lang="ts">
import { toRefs, reactive, onMounted, computed, ref, nextTick, onBeforeUpdate, onBeforeMount, onUnmounted, getCurrentInstance, watch } from 'vue';
import { useRoute, useRouter, onBeforeRouteUpdate } from 'vue-router';
import screenfull from 'screenfull';
import { useStore } from '@/store/index.ts';
import { setSession, getSession, removeSession } from '@/common/utils/storage.ts';
import Sortable from 'sortablejs';
import Contextmenu from '@/views/layout/navBars/tagsView/contextmenu.vue';
export default {
    name: 'layoutTagsView',
    components: { Contextmenu },
    setup() {
        const { proxy } = getCurrentInstance() as any;
        const tagsRefs = ref([]);
        const scrollbarRef = ref();
        const contextmenuRef = ref();
        const tagsUlRef = ref();
        const store = useStore();
        const route = useRoute();
        const router = useRouter();
        const state: any = reactive({
            routePath: route.fullPath,
            dropdown: { x: '', y: '' },
            tagsRefsIndex: 0,
            tagsViewList: [],
            sortable: '',
        });
        // ???????????? tagsView ????????????
        const setTagsStyle = computed(() => {
            return store.state.themeConfig.themeConfig.tagsStyle;
        });
        // ????????????????????????
        const getThemeConfig = computed(() => {
            return store.state.themeConfig.themeConfig;
        });
        // ?????? tagsViewList ????????????????????????????????????????????????????????????
        const addBrowserSetSession = (tagsViewList: Array<object>) => {
            setSession('tagsViewList', tagsViewList);
        };
        // ?????? vuex ?????? tagsViewRoutes ??????
        const getTagsViewRoutes = () => {
            state.routePath = route.fullPath;
            state.tagsViewList = [];
            if (!store.state.themeConfig.themeConfig.isCacheTagsView) removeSession('tagsViewList');
            initTagsView();
        };
        // vuex ??????????????????????????????????????????????????????isAffix???????????????????????????
        const initTagsView = () => {
            if (getSession('tagsViewList') && store.state.themeConfig.themeConfig.isCacheTagsView) {
                state.tagsViewList = getSession('tagsViewList');
            } else {
                // state.tagsViews.map((v: any) => {
                // 	if (v.meta.isAffix && !v.meta.isHide) state.tagsViewList.push({ ...v });
                // });
                addTagsView(route.fullPath);
            }
            // ?????????????????????(li)?????????
            getTagsRefsIndex(route.fullPath);
            // ???????????????????????????????????????????????????
            tagsViewmoveToCurrentTag();
        };

        // 1????????? tagsView?????????????????????isHide?????????????????? tagsView ???
        // path???fullPath
        const addTagsView = (path: string, to: any = null) => {
            if (!to) {
                to = route;
            }
            path = decodeURI(path);
            for (let tv of state.tagsViewList) {
                if (tv.fullPath === path) {
                    return false;
                }
            }
            state.tagsViewList.push({ ...to });
            // addBrowserSetSession(state.tagsViewList);
        };

        // 2??????????????? tagsView???
        // path???fullPath
        const refreshCurrentTagsView = (path: string) => {
            proxy.mittBus.emit('onTagsViewRefreshRouterView', path);
        };

        // 3??????????????? tagsView?????????????????????????????????isAffix?????????????????????
        // path???fullPath
        const closeCurrentTagsView = (path: string) => {
            console.log(path)
            state.tagsViewList.map((v: any, k: number, arr: any) => {
                if (!v.meta.isAffix) {
                    if (v.fullPath === path) {
                        state.tagsViewList.splice(k, 1);
                        setTimeout(() => {
                            // ????????????
                            if (state.tagsViewList.length === k) router.push({ path: arr[arr.length - 1].path, query: arr[arr.length - 1].query });
                            // ???????????????????????????
                            else router.push({ path: arr[k].path, query: arr[k].query });
                        }, 0);
                    }
                }
            });
            // addBrowserSetSession(state.tagsViewList);
        };

        // 4??????????????? tagsView?????????????????????????????????isAffix?????????????????????
        const closeOtherTagsView = (path: string) => {
            const oldTagViews = state.tagsViewList;
            state.tagsViewList = [];
            oldTagViews.map((v: any) => {
                if (v.meta.isAffix && !v.meta.isHide) state.tagsViewList.push({ ...v });
            });
            addTagsView(path);
        };

        // 5??????????????? tagsView?????????????????????????????????isAffix?????????????????????
        const closeAllTagsView = (path: string) => {
            const oldTagViews = state.tagsViewList;
            state.tagsViewList = [];
            oldTagViews.map((v: any) => {
                if (v.meta.isAffix && !v.meta.isHide) {
                    state.tagsViewList.push({ ...v });
                    if (state.tagsViewList.some((v: any) => v.path === path)) router.push({ path, query: route.query });
                    else router.push({ path: v.path, query: route.query });
                }
            });
            // addBrowserSetSession(state.tagsViewList);
        };
        // 6???????????????????????????
        const openCurrenFullscreen = (path: string) => {
            const item = state.tagsViewList.find((v: any) => v.fullPath === path);
            nextTick(() => {
                router.push({ path, query: item.query });
                const element = document.querySelector('.layout-main');
                const screenfulls: any = screenfull;
                screenfulls.request(element);
            });
        };
        // ???????????????????????????
        const onCurrentContextmenuClick = (data: any) => {
            // path???fullPath
            let { id, path } = data;
            let currentTag = state.tagsViewList.find((v: any) => v.fullPath === path);
            switch (id) {
                case 0:
                    refreshCurrentTagsView(path);
                    router.push({ path, query: currentTag.query });
                    break;
                case 1:
                    closeCurrentTagsView(path);
                    break;
                case 2:
                    router.push({ path, query: currentTag.query });
                    closeOtherTagsView(path);
                    break;
                case 3:
                    closeAllTagsView(path);
                    break;
                case 4:
                    openCurrenFullscreen(path);
                    break;
            }
        };
        // ??????????????????
        const isActive = (route: any) => {
            return route.fullPath === state.routePath;
        };
        // ????????????????????? x,y ???????????????????????????props???
        const onContextmenu = (v: any, e: any) => {
            const { clientX, clientY } = e;
            state.dropdown.x = clientX;
            state.dropdown.y = clientY;
            contextmenuRef.value.openContextmenu(v);
        };
        // ????????? tagsView ????????????
        const onTagsClick = (v: any, k: number) => {
            state.routePath = decodeURI(v.fullPath);
            state.tagsRefsIndex = k;
            router.push(v);
        };
        // ?????????????????????
        const updateScrollbar = () => {
            proxy.$refs.scrollbarRef.update();
        };
        // ??????????????????
        const onHandleScroll = (e: any) => {
            proxy.$refs.scrollbarRef.$refs.wrapRef.scrollLeft += e.wheelDelta / 4;
        };
        // tagsView ????????????
        const tagsViewmoveToCurrentTag = () => {
            nextTick(() => {
                if (tagsRefs.value.length <= 0) return false;
                // ?????? li ??????
                let liDom = tagsRefs.value[state.tagsRefsIndex];
                // ?????? li ????????????
                let liIndex = state.tagsRefsIndex;
                // ?????? ul ??? li ???????????????
                let liLength = tagsRefs.value.length;
                // ?????? li
                let liFirst: any = tagsRefs.value[0];
                // ?????? li
                let liLast: any = tagsRefs.value[tagsRefs.value.length - 1];
                // ?????????????????????
                let scrollRefs = proxy.$refs.scrollbarRef.$refs.wrapRef;
                // ???????????????????????????
                let scrollS = scrollRefs.scrollWidth;
                // ???????????????????????????
                let offsetW = scrollRefs.offsetWidth;
                // ???????????????????????????
                let scrollL = scrollRefs.scrollLeft;
                // ????????? tags li dom
                let liPrevTag: any = tagsRefs.value[state.tagsRefsIndex - 1];
                // ????????? tags li dom
                let liNextTag: any = tagsRefs.value[state.tagsRefsIndex + 1];
                // ????????? tags li dom ???????????????
                let beforePrevL: any = '';
                // ????????? tags li dom ???????????????
                let afterNextL: any = '';
                if (liDom === liFirst) {
                    // ??????
                    scrollRefs.scrollLeft = 0;
                } else if (liDom === liLast) {
                    // ??????
                    scrollRefs.scrollLeft = scrollS - offsetW;
                } else {
                    // ??????/??????
                    if (liIndex === 0) beforePrevL = liFirst.offsetLeft - 5;
                    else beforePrevL = liPrevTag?.offsetLeft - 5;
                    if (liIndex === liLength) afterNextL = liLast.offsetLeft + liLast.offsetWidth + 5;
                    else afterNextL = liNextTag.offsetLeft + liNextTag.offsetWidth + 5;
                    if (afterNextL > scrollL + offsetW) {
                        scrollRefs.scrollLeft = afterNextL - offsetW;
                    } else if (beforePrevL < scrollL) {
                        scrollRefs.scrollLeft = beforePrevL;
                    }
                }
                // ?????????????????????????????????
                updateScrollbar();
            });
        };
        // ?????? tagsView ???????????????????????? tagsView ????????????????????????
        const getTagsRefsIndex = (path: string) => {
            if (state.tagsViewList.length > 0) {
                state.tagsRefsIndex = state.tagsViewList.findIndex((item: any) => item.fullPath === path);
            }
        };
        // ?????? tagsView ??????????????????
        const initSortable = () => {
            const el: any = document.querySelector('.layout-navbars-tagsview-ul');
            if (!el) return false;
            if (!getThemeConfig.value.isSortableTagsView) state.sortable && state.sortable.destroy();
            if (getThemeConfig.value.isSortableTagsView) {
                state.sortable = Sortable.create(el, {
                    animation: 300,
                    dataIdAttr: 'data-name',
                    onEnd: () => {
                        const sortEndList: any = [];
                        state.sortable.toArray().map((val: any) => {
                            state.tagsViewList.map((v: any) => {
                                if (v.name === val) sortEndList.push({ ...v });
                            });
                        });
                        // addBrowserSetSession(sortEndList);
                    },
                });
            }
        };
        // ??????????????????????????????????????? tagsView
        // watch(store.state, (val) => {
        // 	if (val.tagsViews.tagsViews.length === state.tagsViewList.length) return false;
        // 	getTagsViewRoutes();
        // });
        // ???????????????
        onBeforeMount(() => {
            // ???????????????????????? 0 ???????????????1 ???????????????2 ???????????????3 ???????????? 4 ???????????????
            proxy.mittBus.on('onCurrentContextmenuClick', (data: object) => {
                onCurrentContextmenuClick(data);
            });
            // ??????????????????????????????/????????????
            proxy.mittBus.on('openOrCloseSortable', () => {
                initSortable();
            });
        });
        // ???????????????
        onUnmounted(() => {
            // ??????????????????????????????
            proxy.mittBus.off('onCurrentContextmenuClick');
            // ????????????????????????????????????/????????????
            proxy.mittBus.off('openOrCloseSortable');
        });
        // ???????????????
        onBeforeUpdate(() => {
            tagsRefs.value = [];
        });
        // ???????????????
        onMounted(() => {
            // ????????? tagsViewRoutes ??????
            getTagsViewRoutes();
            initSortable();
        });
        // ???????????????
        onBeforeRouteUpdate((to) => {
            state.routePath = decodeURI(to.fullPath);
            addTagsView(to.fullPath, to);
            getTagsRefsIndex(to.fullPath);
            tagsViewmoveToCurrentTag();
        });
        return {
            isActive,
            onContextmenu,
            getTagsViewRoutes,
            onTagsClick,
            tagsRefs,
            contextmenuRef,
            scrollbarRef,
            tagsUlRef,
            onHandleScroll,
            getThemeConfig,
            setTagsStyle,
            refreshCurrentTagsView,
            closeCurrentTagsView,
            onCurrentContextmenuClick,
            ...toRefs(state),
        };
    },
};
</script>

<style scoped lang="scss">
.layout-navbars-tagsview {
    flex: 1;
    background-color: #ffffff;
    border-bottom: 1px solid #f1f2f3;
    ::v-deep(.el-scrollbar__wrap) {
        overflow-x: auto !important;
    }
    &-ul {
        list-style: none;
        margin: 0;
        padding: 0;
        height: 34px;
        display: flex;
        align-items: center;
        color: #606266;
        font-size: 12px;
        white-space: nowrap;
        padding: 0 15px;
        &-li {
            height: 26px;
            line-height: 26px;
            display: flex;
            align-items: center;
            border: 1px solid #e6e6e6;
            padding: 0 15px;
            margin-right: 5px;
            border-radius: 2px;
            position: relative;
            z-index: 0;
            cursor: pointer;
            justify-content: space-between;
            &:hover {
                background-color: var(--color-primary-light-9);
                color: var(--color-primary);
                border-color: var(--color-primary-light-6);
            }
            &-iconfont {
                position: relative;
                left: -5px;
                font-size: 12px;
            }
            &-icon {
                border-radius: 100%;
                position: relative;
                height: 14px;
                width: 14px;
                text-align: center;
                line-height: 14px;
                right: -5px;
                &:hover {
                    color: #fff;
                    background-color: var(--color-primary-light-3);
                }
            }
            .layout-icon-active {
                display: block;
            }
            .layout-icon-three {
                display: none;
            }
        }
        .is-active {
            color: #ffffff;
            background: var(--color-primary);
            border-color: var(--color-primary);
        }
    }
    // ??????2
    .tags-style-two {
        .layout-navbars-tagsview-ul-li {
            height: 34px !important;
            line-height: 34px !important;
            border: none !important;
            .layout-navbars-tagsview-ul-li-iconfont {
                display: none;
            }
            .layout-icon-active {
                display: none;
            }
            .layout-icon-three {
                display: block;
            }
        }
        .is-active {
            background: none !important;
            color: var(--color-primary) !important;
            border-bottom: 2px solid !important;
            border-color: var(--color-primary) !important;
            border-radius: 0 !important;
        }
    }
    // ??????3
    .tags-style-three {
        .layout-navbars-tagsview-ul-li {
            height: 34px !important;
            line-height: 34px !important;
            border-right: 1px solid #f6f6f6 !important;
            border-top: none !important;
            border-bottom: none !important;
            border-left: none !important;
            border-radius: 0 !important;
            margin-right: 0 !important;
            &:first-of-type {
                border-left: 1px solid #f6f6f6 !important;
            }
            .layout-icon-active {
                display: none;
            }
            .layout-icon-three {
                display: block;
            }
        }
        .is-active {
            background: white !important;
            color: var(--color-primary) !important;
            border-top: 1px solid !important;
            border-top-color: var(--color-primary) !important;
        }
    }
    // ??????4
    .tags-style-four {
        .layout-navbars-tagsview-ul-li {
            margin-right: 0 !important;
            border: none !important;
            position: relative;
            border-radius: 3px !important;
            .layout-icon-active {
                display: none;
            }
            .layout-icon-three {
                display: block;
            }
            &:hover {
                background: none !important;
            }
        }
        .is-active {
            background: none !important;
            color: var(--color-primary) !important;
        }
    }
}
.layout-navbars-tagsview-shadow {
    box-shadow: rgb(0 21 41 / 4%) 0px 1px 4px;
}
</style>
