import jsYaml from 'js-yaml';

class PkgFetcher {
    async fetch(pkgRef) {
        return fetch(`/pkgs/${encodeURIComponent(pkgRef)}/contents/op.yml`)
            .then(response => (response.text()))
            .then(PkgFetcher._parse);
    }

    static _parse(opDotYml){
        return jsYaml.safeLoad(opDotYml);
    }
}

// export as singleton
export default new PkgFetcher();
