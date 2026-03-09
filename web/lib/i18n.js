(function (global) {
    const DEFAULT_ALIASES = {
        zh: 'zh-CN',
        'zh-cn': 'zh-CN',
        'zh-sg': 'zh-CN',
        'zh-hans': 'zh-CN',
        'zh-tw': 'zh-CN',
        'zh-hk': 'zh-CN',
        en: 'en',
        'en-us': 'en',
        'en-gb': 'en',
        'en-au': 'en',
        ja: 'ja',
        'ja-jp': 'ja',
        es: 'es',
        'es-es': 'es',
        'es-mx': 'es',
        'es-419': 'es'
    };

    function interpolate(template, vars) {
        return String(template).replace(/\{(\w+)\}/g, function (_, key) {
            return Object.prototype.hasOwnProperty.call(vars, key) ? vars[key] : `{${key}}`;
        });
    }

    function getValue(source, path) {
        return String(path || '')
            .split('.')
            .filter(Boolean)
            .reduce(function (current, key) {
                if (current == null) {
                    return undefined;
                }
                return current[key];
            }, source);
    }

    function findSupportedLanguage(translations, requested) {
        const supported = Object.keys(translations || {});
        const normalizedRequested = String(requested || '').trim().replace(/_/g, '-').toLowerCase();

        if (!normalizedRequested) {
            return null;
        }

        const exact = supported.find(function (candidate) {
            return candidate.toLowerCase() === normalizedRequested;
        });
        if (exact) {
            return exact;
        }

        const aliased = DEFAULT_ALIASES[normalizedRequested];
        if (aliased && translations[aliased]) {
            return aliased;
        }

        const base = normalizedRequested.split('-')[0];
        const baseAlias = DEFAULT_ALIASES[base];
        if (baseAlias && translations[baseAlias]) {
            return baseAlias;
        }

        const baseMatch = supported.find(function (candidate) {
            return candidate.toLowerCase().split('-')[0] === base;
        });
        return baseMatch || null;
    }

    function resolveLanguage(translations, fallback) {
        const candidates = [];

        if (Array.isArray(global.navigator && global.navigator.languages)) {
            candidates.push.apply(candidates, global.navigator.languages);
        }

        if (global.navigator && global.navigator.language) {
            candidates.push(global.navigator.language);
        }

        for (let index = 0; index < candidates.length; index += 1) {
            const match = findSupportedLanguage(translations, candidates[index]);
            if (match) {
                return match;
            }
        }

        if (translations && translations[fallback]) {
            return fallback;
        }

        return Object.keys(translations || {})[0] || fallback || 'en';
    }

    function create(options) {
        const translations = options && options.translations ? options.translations : {};
        const fallback = options && options.fallback ? options.fallback : Object.keys(translations)[0] || 'en';
        const language = resolveLanguage(translations, fallback);
        const currentMessages = translations[language] || translations[fallback] || {};
        const fallbackMessages = translations[fallback] || {};

        function t(path, vars) {
            const value = getValue(currentMessages, path);
            const resolved = value === undefined ? getValue(fallbackMessages, path) : value;

            if (typeof resolved === 'function') {
                return resolved(vars || {});
            }

            if (typeof resolved === 'string') {
                return interpolate(resolved, vars || {});
            }

            return resolved !== undefined ? resolved : '';
        }

        function applyMeta() {
            const lang = t('meta.lang') || language;
            if (lang) {
                document.documentElement.lang = lang;
            }

            const title = t('meta.title');
            if (title) {
                document.title = title;
            }

            const description = document.querySelector('meta[name="description"]');
            if (description && t('meta.description')) {
                description.setAttribute('content', t('meta.description'));
            }

            const keywords = document.querySelector('meta[name="keywords"]');
            if (keywords && t('meta.keywords')) {
                keywords.setAttribute('content', t('meta.keywords'));
            }
        }

        function selectAll(target) {
            if (!target) {
                return [];
            }

            if (typeof target === 'string') {
                return Array.from(document.querySelectorAll(target));
            }

            if (target instanceof Element) {
                return [target];
            }

            if (Array.isArray(target)) {
                return target.flatMap(selectAll);
            }

            return [];
        }

        function apply(definitions) {
            (definitions || []).forEach(function (definition) {
                const elements = selectAll(definition.selector);
                if (!elements.length) {
                    return;
                }

                elements.forEach(function (element, index) {
                    let value = definition.value;
                    if (value === undefined && definition.key) {
                        const vars = typeof definition.vars === 'function'
                            ? definition.vars(element, index)
                            : definition.vars;
                        value = t(definition.key, vars || {});
                    }

                    if (typeof definition.transform === 'function') {
                        value = definition.transform(value, element, index);
                    }

                    if (value === undefined || value === null) {
                        return;
                    }

                    if (definition.attr) {
                        element.setAttribute(definition.attr, value);
                    } else if (definition.html) {
                        element.innerHTML = value;
                    } else {
                        element.textContent = value;
                    }
                });
            });
        }

        return {
            language,
            fallback,
            t,
            apply,
            applyMeta,
            messages: function () {
                return currentMessages;
            }
        };
    }

    global.CliGoolI18n = {
        create: create
    };
})(window);
