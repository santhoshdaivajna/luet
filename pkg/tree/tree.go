// Recipe is a builder imeplementation.

// It reads a Tree and spit it in human readable form (YAML), called recipe,
// It also loads a tree (recipe) from a YAML (to a db, e.g. BoltDB), allowing to query it
// with the solver, using the package object.
package tree

import (
	"errors"
	"fmt"

	pkg "github.com/mudler/luet/pkg/package"
)

func NewDefaultTree() pkg.Tree { return &DefaultTree{} }

type DefaultTree struct {
	Packages   pkg.PackageSet
	CacheWorld []pkg.Package
}

func (gt *DefaultTree) GetPackageSet() pkg.PackageSet {
	return gt.Packages
}

func (gt *DefaultTree) Prelude() string {
	return ""
}

func (gt *DefaultTree) SetPackageSet(s pkg.PackageSet) {
	gt.Packages = s
}

func (gt *DefaultTree) World() ([]pkg.Package, error) {
	if len(gt.CacheWorld) > 0 {
		return gt.CacheWorld, nil
	}
	packages := []pkg.Package{}
	for _, pid := range gt.GetPackageSet().GetPackages() {

		p, err := gt.GetPackageSet().GetPackage(pid)
		if err != nil {
			return packages, err
		}
		packages = append(packages, p)
	}
	gt.CacheWorld = packages
	return packages, nil
}

// FIXME: Dup in Packageset
func (gt *DefaultTree) FindPackage(pack pkg.Package) (pkg.Package, error) {
	packages, err := gt.World()
	if err != nil {
		return nil, err
	}
	for _, pid := range packages {
		if pack.GetFingerPrint() == pid.GetFingerPrint() {
			return pid, nil
		}
	}
	return nil, errors.New("No package found")
}

// Search for deps/conflicts in db and replaces it with packages in the db
func (t *DefaultTree) ResolveDeps() error {
	for _, pid := range t.GetPackageSet().GetPackages() {

		p, err := t.GetPackageSet().GetPackage(pid)
		if err != nil {
			return err
		}

		for _, r := range p.GetRequires() {

			foundPackage, err := t.GetPackageSet().FindPackage(r)
			if err != nil {
				fmt.Println("Warning: Unmatched dependency - no package found in the database for this requirement clause")
				continue
				//return err
			}
			found, ok := foundPackage.(*pkg.DefaultPackage)
			if !ok {
				panic("Simpleparser should deal only with DefaultPackages")
			}
			r = found
		}

		for _, r := range p.GetConflicts() {

			foundPackage, err := t.GetPackageSet().FindPackage(r)
			if err != nil {
				continue
				//return err
			}
			found, ok := foundPackage.(*pkg.DefaultPackage)
			if !ok {
				panic("Simpleparser should deal only with DefaultPackages")
			}
			r = found
		}

		if err = t.GetPackageSet().UpdatePackage(p); err != nil {
			return err
		}

	}
	return nil
}