// Copyright © 2019 Ettore Di Giacinto <mudler@gentoo.org>
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, see <http://www.gnu.org/licenses/>.

package pkg_test

import (
	"github.com/crillab/gophersat/bf"
	. "github.com/mudler/luet/pkg/package"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Package", func() {

	Context("Simple package", func() {
		a := NewPackage("A", ">=1.0", []*DefaultPackage{}, []*DefaultPackage{})
		a1 := NewPackage("A", "1.0", []*DefaultPackage{}, []*DefaultPackage{})
		a11 := NewPackage("A", "1.1", []*DefaultPackage{}, []*DefaultPackage{})
		a01 := NewPackage("A", "0.1", []*DefaultPackage{}, []*DefaultPackage{})
		It("Expands correctly", func() {
			lst, err := a.Expand([]Package{a1, a11, a01})
			Expect(err).ToNot(HaveOccurred())
			Expect(lst).To(ContainElement(a11))
			Expect(lst).To(ContainElement(a1))
			Expect(lst).ToNot(ContainElement(a01))
			Expect(len(lst)).To(Equal(2))
		})
	})

	Context("RequiresContains", func() {
		a := NewPackage("A", ">=1.0", []*DefaultPackage{}, []*DefaultPackage{})
		a1 := NewPackage("A", "1.0", []*DefaultPackage{a}, []*DefaultPackage{})
		a11 := NewPackage("A", "1.1", []*DefaultPackage{}, []*DefaultPackage{})
		a01 := NewPackage("A", "0.1", []*DefaultPackage{a1, a11}, []*DefaultPackage{})
		It("returns correctly", func() {
			Expect(a01.RequiresContains(a1)).To(BeTrue())
			Expect(a01.RequiresContains(a11)).To(BeTrue())
			Expect(a01.RequiresContains(a)).To(BeTrue())
			Expect(a.RequiresContains(a11)).ToNot(BeTrue())
		})
	})

	Context("Encoding", func() {
		a1 := NewPackage("A", "1.0", []*DefaultPackage{}, []*DefaultPackage{})
		a11 := NewPackage("A", "1.1", []*DefaultPackage{}, []*DefaultPackage{})
		a := NewPackage("A", ">=1.0", []*DefaultPackage{a1}, []*DefaultPackage{a11})
		It("decodes and encodes correctly", func() {

			ID, err := a.Encode()
			Expect(err).ToNot(HaveOccurred())

			p, err := DecodePackage(ID)
			Expect(err).ToNot(HaveOccurred())

			Expect(p.GetVersion()).To(Equal(a.GetVersion()))
			Expect(p.GetName()).To(Equal(a.GetName()))
			Expect(p.Flagged()).To(Equal(a.Flagged()))
			Expect(p.GetFingerPrint()).To(Equal(a.GetFingerPrint()))
			Expect(len(p.GetConflicts())).To(Equal(len(a.GetConflicts())))
			Expect(len(p.GetRequires())).To(Equal(len(a.GetRequires())))
			Expect(len(p.GetRequires())).To(Equal(1))
			Expect(len(p.GetConflicts())).To(Equal(1))
			Expect(p.GetConflicts()[0].GetName()).To(Equal(a11.GetName()))
			Expect(a.GetConflicts()[0].GetName()).To(Equal(a11.GetName()))
			Expect(p.GetRequires()[0].GetName()).To(Equal(a1.GetName()))
			Expect(a.GetRequires()[0].GetName()).To(Equal(a1.GetName()))
		})
	})

	Context("BuildFormula", func() {
		It("builds empty constraints", func() {
			a1 := NewPackage("A", "1.0", []*DefaultPackage{}, []*DefaultPackage{})
			f, err := a1.BuildFormula()
			Expect(err).ToNot(HaveOccurred())
			Expect(f).To(Equal([]bf.Formula(nil)))
		})
		It("builds constraints correctly", func() {
			a11 := NewPackage("A", "1.1", []*DefaultPackage{}, []*DefaultPackage{})
			a21 := NewPackage("A", "1.2", []*DefaultPackage{}, []*DefaultPackage{})
			a1 := NewPackage("A", "1.0", []*DefaultPackage{a11}, []*DefaultPackage{a21})
			f, err := a1.BuildFormula()
			Expect(err).ToNot(HaveOccurred())
			Expect(len(f)).To(Equal(2))
			Expect(f[0].String()).To(Equal("or(not(1c6f6460), 0e78200a)"))
			Expect(f[1].String()).To(Equal("or(not(1c6f6460), not(039f5fc3))"))
		})
	})
})
